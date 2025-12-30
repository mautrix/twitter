package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000D\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\f\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\b\b\u0086\b\u0018\u0000 &2\u00020\u0001:\u0002'&B5\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\u000e\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u0006\u0012\b\u0010\t\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\n\u0010\u000bJ\u0017\u0010\u000f\u001a\u00020\u000e2\u0006\u0010\r\u001a\u00020\fH\u0016¢\u0006\u0004\b\u000f\u0010\u0010J\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u0012J\u0012\u0010\u0013\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0013\u0010\u0014J\u0018\u0010\u0015\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0016J\u0012\u0010\u0017\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0017\u0010\u0012JF\u0010\u0018\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\u0010\b\u0002\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u00062\n\b\u0002\u0010\t\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u0018\u0010\u0019J\u0010\u0010\u001a\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u001a\u0010\u0012J\u0010\u0010\u001c\u001a\u00020\u001bHÖ\u0001¢\u0006\u0004\b\u001c\u0010\u001dJ\u001a\u0010!\u001a\u00020 2\b\u0010\u001f\u001a\u0004\u0018\u00010\u001eHÖ\u0003¢\u0006\u0004\b!\u0010\"R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010#R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010$R\u001c\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010%R\u0016\u0010\t\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\t\u0010#¨\u0006("}, m64930d2 = {"Lcom/x/dmv2/thriftjava/KeyRotation;", "Lcom/bendb/thrifty/a;", "", "previous_version", "Lcom/x/dmv2/thriftjava/RatchetTree;", "ratchet_tree", "", "Lcom/x/dmv2/thriftjava/UpdatePathNode;", "nodes", "encrypted_private_key", "<init>", "(Ljava/lang/String;Lcom/x/dmv2/thriftjava/RatchetTree;Ljava/util/List;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Lcom/x/dmv2/thriftjava/RatchetTree;", "component3", "()Ljava/util/List;", "component4", "copy", "(Ljava/lang/String;Lcom/x/dmv2/thriftjava/RatchetTree;Ljava/util/List;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/KeyRotation;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Lcom/x/dmv2/thriftjava/RatchetTree;", "Ljava/util/List;", "Companion", "KeyRotationAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class KeyRotation implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String encrypted_private_key;

    @JvmField
    @InterfaceC88465b
    public final List nodes;

    @JvmField
    @InterfaceC88465b
    public final String previous_version;

    @JvmField
    @InterfaceC88465b
    public final RatchetTree ratchet_tree;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new KeyRotationAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/KeyRotation$KeyRotationAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/KeyRotation;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/KeyRotation;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/KeyRotation;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class KeyRotationAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public KeyRotation m85578read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            RatchetTree ratchetTree = null;
            ArrayList arrayList = null;
            String string2 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new KeyRotation(string, ratchetTree, arrayList, string2);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 4) {
                                C11272a.m14141a(protocol, b);
                            } else if (b == 11) {
                                string2 = protocol.readString();
                            } else {
                                C11272a.m14141a(protocol, b);
                            }
                        } else if (b == 15) {
                            int i = protocol.mo14130a2().f38395b;
                            ArrayList arrayList2 = new ArrayList(i);
                            for (int i2 = 0; i2 < i; i2++) {
                                arrayList2.add((UpdatePathNode) UpdatePathNode.ADAPTER.read(protocol));
                            }
                            arrayList = arrayList2;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 12) {
                        ratchetTree = (RatchetTree) RatchetTree.ADAPTER.read(protocol);
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 11) {
                    string = protocol.readString();
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a KeyRotation struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("KeyRotation");
            if (struct.previous_version != null) {
                protocol.mo14136v3("previous_version", 1, (byte) 11);
                protocol.mo14137w0(struct.previous_version);
            }
            if (struct.ratchet_tree != null) {
                protocol.mo14136v3("ratchet_tree", 2, (byte) 12);
                RatchetTree.ADAPTER.write(protocol, struct.ratchet_tree);
            }
            if (struct.nodes != null) {
                protocol.mo14136v3("nodes", 3, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.nodes.size());
                Iterator it = struct.nodes.iterator();
                while (it.hasNext()) {
                    UpdatePathNode.ADAPTER.write(protocol, (UpdatePathNode) it.next());
                }
            }
            if (struct.encrypted_private_key != null) {
                protocol.mo14136v3("encrypted_private_key", 4, (byte) 11);
                protocol.mo14137w0(struct.encrypted_private_key);
            }
            protocol.mo14134i0();
        }
    }

    public KeyRotation(@InterfaceC88465b String str, @InterfaceC88465b RatchetTree ratchetTree, @InterfaceC88465b List list, @InterfaceC88465b String str2) {
        this.previous_version = str;
        this.ratchet_tree = ratchetTree;
        this.nodes = list;
        this.encrypted_private_key = str2;
    }

    public static /* synthetic */ KeyRotation copy$default(KeyRotation keyRotation, String str, RatchetTree ratchetTree, List list, String str2, int i, Object obj) {
        if ((i & 1) != 0) {
            str = keyRotation.previous_version;
        }
        if ((i & 2) != 0) {
            ratchetTree = keyRotation.ratchet_tree;
        }
        if ((i & 4) != 0) {
            list = keyRotation.nodes;
        }
        if ((i & 8) != 0) {
            str2 = keyRotation.encrypted_private_key;
        }
        return keyRotation.copy(str, ratchetTree, list, str2);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getPrevious_version() {
        return this.previous_version;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final RatchetTree getRatchet_tree() {
        return this.ratchet_tree;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final List getNodes() {
        return this.nodes;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getEncrypted_private_key() {
        return this.encrypted_private_key;
    }

    @InterfaceC88464a
    public final KeyRotation copy(@InterfaceC88465b String previous_version, @InterfaceC88465b RatchetTree ratchet_tree, @InterfaceC88465b List nodes, @InterfaceC88465b String encrypted_private_key) {
        return new KeyRotation(previous_version, ratchet_tree, nodes, encrypted_private_key);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof KeyRotation)) {
            return false;
        }
        KeyRotation keyRotation = (KeyRotation) other;
        return Intrinsics.m65267c(this.previous_version, keyRotation.previous_version) && Intrinsics.m65267c(this.ratchet_tree, keyRotation.ratchet_tree) && Intrinsics.m65267c(this.nodes, keyRotation.nodes) && Intrinsics.m65267c(this.encrypted_private_key, keyRotation.encrypted_private_key);
    }

    public int hashCode() {
        String str = this.previous_version;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        RatchetTree ratchetTree = this.ratchet_tree;
        int iHashCode2 = (iHashCode + (ratchetTree == null ? 0 : ratchetTree.hashCode())) * 31;
        List list = this.nodes;
        int iHashCode3 = (iHashCode2 + (list == null ? 0 : list.hashCode())) * 31;
        String str2 = this.encrypted_private_key;
        return iHashCode3 + (str2 != null ? str2.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "KeyRotation(previous_version=" + this.previous_version + ", ratchet_tree=" + this.ratchet_tree + ", nodes=" + this.nodes + ", encrypted_private_key=" + this.encrypted_private_key + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
