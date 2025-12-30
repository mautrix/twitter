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

@Metadata(m64929d1 = {"\u0000F\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\b\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 !2\u00020\u0001:\u0002\"!B!\u0012\u000e\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0005¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0018\u0010\u000e\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0005HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J.\u0010\u0012\u001a\u00020\u00002\u0010\b\u0002\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u0005HÆ\u0001¢\u0006\u0004\b\u0012\u0010\u0013J\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u0010\u0010\u0018\u001a\u00020\u0017HÖ\u0001¢\u0006\u0004\b\u0018\u0010\u0019J\u001a\u0010\u001d\u001a\u00020\u001c2\b\u0010\u001b\u001a\u0004\u0018\u00010\u001aHÖ\u0003¢\u0006\u0004\b\u001d\u0010\u001eR\u001c\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001fR\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00058\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010 ¨\u0006#"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/StoredGroupState;", "Lcom/bendb/thrifty/a;", "", "Lcom/x/dmv2/thriftjava/MaybeKeypair;", "keypairs", "Lcom/x/dmv2/thriftjava/RatchetTree;", "ratchet_tree", "<init>", "(Ljava/util/List;Lcom/x/dmv2/thriftjava/RatchetTree;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/util/List;", "component2", "()Lcom/x/dmv2/thriftjava/RatchetTree;", "copy", "(Ljava/util/List;Lcom/x/dmv2/thriftjava/RatchetTree;)Lcom/x/dmv2/thriftjava/StoredGroupState;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/util/List;", "Lcom/x/dmv2/thriftjava/RatchetTree;", "Companion", "StoredGroupStateAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class StoredGroupState implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List keypairs;

    @JvmField
    @InterfaceC88465b
    public final RatchetTree ratchet_tree;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new StoredGroupStateAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/StoredGroupState$StoredGroupStateAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/StoredGroupState;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/StoredGroupState;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/StoredGroupState;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class StoredGroupStateAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public StoredGroupState m85591read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            ArrayList arrayList = null;
            RatchetTree ratchetTree = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new StoredGroupState(arrayList, ratchetTree);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 12) {
                        ratchetTree = (RatchetTree) RatchetTree.ADAPTER.read(protocol);
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 15) {
                    int i = protocol.mo14130a2().f38395b;
                    ArrayList arrayList2 = new ArrayList(i);
                    for (int i2 = 0; i2 < i; i2++) {
                        arrayList2.add((MaybeKeypair) MaybeKeypair.ADAPTER.read(protocol));
                    }
                    arrayList = arrayList2;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a StoredGroupState struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("StoredGroupState");
            if (struct.keypairs != null) {
                protocol.mo14136v3("keypairs", 1, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.keypairs.size());
                Iterator it = struct.keypairs.iterator();
                while (it.hasNext()) {
                    MaybeKeypair.ADAPTER.write(protocol, (MaybeKeypair) it.next());
                }
            }
            if (struct.ratchet_tree != null) {
                protocol.mo14136v3("ratchet_tree", 2, (byte) 12);
                RatchetTree.ADAPTER.write(protocol, struct.ratchet_tree);
            }
            protocol.mo14134i0();
        }
    }

    public StoredGroupState(@InterfaceC88465b List list, @InterfaceC88465b RatchetTree ratchetTree) {
        this.keypairs = list;
        this.ratchet_tree = ratchetTree;
    }

    public static /* synthetic */ StoredGroupState copy$default(StoredGroupState storedGroupState, List list, RatchetTree ratchetTree, int i, Object obj) {
        if ((i & 1) != 0) {
            list = storedGroupState.keypairs;
        }
        if ((i & 2) != 0) {
            ratchetTree = storedGroupState.ratchet_tree;
        }
        return storedGroupState.copy(list, ratchetTree);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final List getKeypairs() {
        return this.keypairs;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final RatchetTree getRatchet_tree() {
        return this.ratchet_tree;
    }

    @InterfaceC88464a
    public final StoredGroupState copy(@InterfaceC88465b List keypairs, @InterfaceC88465b RatchetTree ratchet_tree) {
        return new StoredGroupState(keypairs, ratchet_tree);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof StoredGroupState)) {
            return false;
        }
        StoredGroupState storedGroupState = (StoredGroupState) other;
        return Intrinsics.m65267c(this.keypairs, storedGroupState.keypairs) && Intrinsics.m65267c(this.ratchet_tree, storedGroupState.ratchet_tree);
    }

    public int hashCode() {
        List list = this.keypairs;
        int iHashCode = (list == null ? 0 : list.hashCode()) * 31;
        RatchetTree ratchetTree = this.ratchet_tree;
        return iHashCode + (ratchetTree != null ? ratchetTree.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "StoredGroupState(keypairs=" + this.keypairs + ", ratchet_tree=" + this.ratchet_tree + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
