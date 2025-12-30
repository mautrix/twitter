package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000:\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u000e\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u000b\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\t\b\u0086\b\u0018\u0000 \"2\u00020\u0001:\u0002#\"B%\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0006¢\u0006\u0004\b\b\u0010\tJ\u0017\u0010\r\u001a\u00020\f2\u0006\u0010\u000b\u001a\u00020\nH\u0016¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u0012J\u0012\u0010\u0013\u001a\u0004\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u0013\u0010\u0014J4\u0010\u0015\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u0006HÆ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u0010\u0010\u0017\u001a\u00020\u0004HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0012J\u0010\u0010\u0019\u001a\u00020\u0018HÖ\u0001¢\u0006\u0004\b\u0019\u0010\u001aJ\u001a\u0010\u001d\u001a\u00020\u00022\b\u0010\u001c\u001a\u0004\u0018\u00010\u001bHÖ\u0003¢\u0006\u0004\b\u001d\u0010\u001eR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001fR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010 R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u0010!¨\u0006$"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;", "Lcom/bendb/thrifty/a;", "", "finished_pull", "", "sequence_continue", "Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "pull_message_page_details", "<init>", "(Ljava/lang/Boolean;Ljava/lang/String;Lcom/x/dmv2/thriftjava/PullMessagePageDetails;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Boolean;", "component2", "()Ljava/lang/String;", "component3", "()Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "copy", "(Ljava/lang/Boolean;Ljava/lang/String;Lcom/x/dmv2/thriftjava/PullMessagePageDetails;)Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;", "toString", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Boolean;", "Ljava/lang/String;", "Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "Companion", "PullMessagesFinishedInstructionAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class PullMessagesFinishedInstruction implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Boolean finished_pull;

    @JvmField
    @InterfaceC88465b
    public final PullMessagePageDetails pull_message_page_details;

    @JvmField
    @InterfaceC88465b
    public final String sequence_continue;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new PullMessagesFinishedInstructionAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction$PullMessagesFinishedInstructionAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class PullMessagesFinishedInstructionAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public PullMessagesFinishedInstruction m85668read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Boolean boolValueOf = null;
            String string = null;
            PullMessagePageDetails pullMessagePageDetails = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new PullMessagesFinishedInstruction(boolValueOf, string, pullMessagePageDetails);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 12) {
                            pullMessagePageDetails = (PullMessagePageDetails) PullMessagePageDetails.ADAPTER.read(protocol);
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 11) {
                        string = protocol.readString();
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 2) {
                    boolValueOf = Boolean.valueOf(protocol.readBool());
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a PullMessagesFinishedInstruction struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("PullMessagesFinishedInstruction");
            if (struct.finished_pull != null) {
                protocol.mo14136v3("finished_pull", 1, (byte) 2);
                protocol.mo14125P1(struct.finished_pull.booleanValue());
            }
            if (struct.sequence_continue != null) {
                protocol.mo14136v3("sequence_continue", 2, (byte) 11);
                protocol.mo14137w0(struct.sequence_continue);
            }
            if (struct.pull_message_page_details != null) {
                protocol.mo14136v3("pull_message_page_details", 3, (byte) 12);
                PullMessagePageDetails.ADAPTER.write(protocol, struct.pull_message_page_details);
            }
            protocol.mo14134i0();
        }
    }

    public PullMessagesFinishedInstruction(@InterfaceC88465b Boolean bool, @InterfaceC88465b String str, @InterfaceC88465b PullMessagePageDetails pullMessagePageDetails) {
        this.finished_pull = bool;
        this.sequence_continue = str;
        this.pull_message_page_details = pullMessagePageDetails;
    }

    public static /* synthetic */ PullMessagesFinishedInstruction copy$default(PullMessagesFinishedInstruction pullMessagesFinishedInstruction, Boolean bool, String str, PullMessagePageDetails pullMessagePageDetails, int i, Object obj) {
        if ((i & 1) != 0) {
            bool = pullMessagesFinishedInstruction.finished_pull;
        }
        if ((i & 2) != 0) {
            str = pullMessagesFinishedInstruction.sequence_continue;
        }
        if ((i & 4) != 0) {
            pullMessagePageDetails = pullMessagesFinishedInstruction.pull_message_page_details;
        }
        return pullMessagesFinishedInstruction.copy(bool, str, pullMessagePageDetails);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Boolean getFinished_pull() {
        return this.finished_pull;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getSequence_continue() {
        return this.sequence_continue;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final PullMessagePageDetails getPull_message_page_details() {
        return this.pull_message_page_details;
    }

    @InterfaceC88464a
    public final PullMessagesFinishedInstruction copy(@InterfaceC88465b Boolean finished_pull, @InterfaceC88465b String sequence_continue, @InterfaceC88465b PullMessagePageDetails pull_message_page_details) {
        return new PullMessagesFinishedInstruction(finished_pull, sequence_continue, pull_message_page_details);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof PullMessagesFinishedInstruction)) {
            return false;
        }
        PullMessagesFinishedInstruction pullMessagesFinishedInstruction = (PullMessagesFinishedInstruction) other;
        return Intrinsics.m65267c(this.finished_pull, pullMessagesFinishedInstruction.finished_pull) && Intrinsics.m65267c(this.sequence_continue, pullMessagesFinishedInstruction.sequence_continue) && Intrinsics.m65267c(this.pull_message_page_details, pullMessagesFinishedInstruction.pull_message_page_details);
    }

    public int hashCode() {
        Boolean bool = this.finished_pull;
        int iHashCode = (bool == null ? 0 : bool.hashCode()) * 31;
        String str = this.sequence_continue;
        int iHashCode2 = (iHashCode + (str == null ? 0 : str.hashCode())) * 31;
        PullMessagePageDetails pullMessagePageDetails = this.pull_message_page_details;
        return iHashCode2 + (pullMessagePageDetails != null ? pullMessagePageDetails.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "PullMessagesFinishedInstruction(finished_pull=" + this.finished_pull + ", sequence_continue=" + this.sequence_continue + ", pull_message_page_details=" + this.pull_message_page_details + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
